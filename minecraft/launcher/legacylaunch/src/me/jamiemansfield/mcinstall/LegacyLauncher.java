// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package me.jamiemansfield.mcinstall;

import net.minecraft.Launcher;

import javax.swing.JPanel;
import java.applet.Applet;
import java.awt.BorderLayout;
import java.awt.Color;
import java.awt.Dimension;
import java.awt.Frame;
import java.awt.event.WindowAdapter;
import java.awt.event.WindowEvent;
import java.io.File;
import java.lang.invoke.MethodHandle;
import java.lang.invoke.MethodHandles;
import java.lang.invoke.MethodType;
import java.lang.reflect.Field;
import java.lang.reflect.Modifier;
import java.nio.file.Path;
import java.nio.file.Paths;

public final class LegacyLauncher {

    private static final int DEFAULT_WIDTH = 854;
    private static final int DEFAULT_HEIGHT = 480;

    private static final ClassLoader loader = LegacyLauncher.class.getClassLoader();

    public static void main(final String[] args) {
        Thread.currentThread().setUncaughtExceptionHandler((t, e) -> {
            System.out.println("Uncaught exception!");
            e.printStackTrace(System.out);
        });

        // Create frame for wrapped applet
        final Frame frame = createFrame("Tekkit Classic 3.1.2");
        frame.setVisible(true);

        // Give the frame the default initial size, using a temporary panel with a
        // preferred size. Oddly using Frame#setSize(Dimension) caused Minecraft to
        // be stuck in one GUI scale (extremely small) - so for now, despite how ugly
        // this is, it will do.
        final JPanel panel = new JPanel();
        frame.setLayout(new BorderLayout());
        panel.setPreferredSize(new Dimension(DEFAULT_WIDTH, DEFAULT_HEIGHT));
        frame.add(panel, BorderLayout.CENTER);
        frame.pack();
        frame.setVisible(true);

        try {
            // Rewrite the Minecraft class to use the correct game directory
            rewriteMinecraft(Paths.get(""));

            final Launcher launcher = new Launcher(createMinecraftApplet());

            frame.removeAll();
            frame.setLayout(new BorderLayout());
            frame.add(launcher, BorderLayout.CENTER);
            frame.validate();

            launcher.init();
            launcher.start();

            Runtime.getRuntime().addShutdownHook(new Thread(launcher::stop));
        }
        catch (final Throwable ex) {
            ex.printStackTrace(System.out);
            System.exit(-1);
        }
    }

    private static void rewriteMinecraft(final Path dir) throws Throwable {
        // Get the Minecraft class
        Class<?> minecraftClass;
        try {
            minecraftClass = loader.loadClass("net.minecraft.client.Minecraft");
        }
        catch (final ClassNotFoundException ignored) {
            minecraftClass = loader.loadClass("com.mojang.minecraft.Minecraft");
        }
        System.out.println("Using '" + minecraftClass.getName() + "' as Minecraft class");

        final Path minecraftRoot = dir.toAbsolutePath();
        System.out.println("Using '" + minecraftRoot.toString() + "' as Minecraft root");

        // Find root field
        for (final Field field : minecraftClass.getDeclaredFields()) {
            // Check the modifiers are correct
            if (!Modifier.isStatic(field.getModifiers()) || !Modifier.isPrivate(field.getModifiers())) {
                continue;
            }
            // Check the type is correct
            if (field.getType() != File.class) {
                continue;
            }

            System.out.println("Found minecraft root field: " + field.getName());
            field.setAccessible(true);
            MethodHandles.lookup().unreflectSetter(field)
                    .invoke(minecraftRoot.toFile());

            return;
        }

        // Failed to find Minecraft root field, crash out
        throw new RuntimeException("Failed to find minecraft root field");
    }

    private static Applet createMinecraftApplet() throws Throwable {
        // Get the Minecraft applet class
        Class<?> appletClass;
        try {
            appletClass = loader.loadClass("net.minecraft.client.MinecraftApplet");
        }
        catch (final ClassNotFoundException ignored) {
            appletClass = loader.loadClass("com.mojang.minecraft.MinecraftApplet");
        }
        System.out.println("Using '" + appletClass.getName() + "' as applet class");

        // Create an instance
        final MethodType descriptor = MethodType.methodType(void.class);
        final MethodHandle handle = MethodHandles.lookup().findConstructor(appletClass, descriptor);
        return (Applet) handle.invoke();
    }

    private static Frame createFrame(final String title) {
        System.out.println("Creating frame with title '" + title + "'");

        // I found that using JFrame didn't work, but Frame did - so that's why this
        // isn't a Swing frame :S
        final Frame frame = new Frame(title);
        frame.setBackground(Color.BLACK);
        frame.setLocationRelativeTo(null);
        frame.addWindowListener(new WindowAdapter() {
            @Override
            public void windowClosing(final WindowEvent e) {
                System.exit(1);
            }
        });
        return frame;
    }

    private LegacyLauncher() {
    }

}
