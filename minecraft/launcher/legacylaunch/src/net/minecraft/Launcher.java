// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package net.minecraft;

import java.applet.Applet;
import java.applet.AppletStub;
import java.awt.BorderLayout;
import java.net.MalformedURLException;
import java.net.URL;
import java.util.Map;

public class Launcher extends Applet implements AppletStub {

    private static final String BASE = "http://www.minecraft.net/game/";

    private final Applet minecraftApplet;
    private final Map<String, String> args;

    public Launcher(final Applet minecraftApplet, final Map<String, String> args) {
        this.minecraftApplet = minecraftApplet;
        this.minecraftApplet.setStub(this);
        this.args = args;

        this.setLayout(new BorderLayout());
        this.add(minecraftApplet, BorderLayout.CENTER);
        this.validate();
    }

    @Override
    public boolean isActive() {
        return true;
    }

    @Override
    public URL getDocumentBase() {
        try {
            return new URL(BASE);
        }
        catch (final MalformedURLException ignored) {
        }
        return null;
    }

    @Override
    public URL getCodeBase() {
        try {
            return new URL(BASE);
        }
        catch (final MalformedURLException ignored) {
        }
        return null;
    }

    @Override
    public String getParameter(final String name) {
        final String value = this.args.get(name);
        if (value != null) return value;

        System.out.println("Client asked for '" + name + "' parameter");

        switch (name) {
            // Show the quit button
            case "stand-alone":
            // Allow players to save their levels
            case "haspaid":
                return "true";
            case "demo":
            case "fullscreen":
                return "false";
        }

        return null;
    }

    @Override
    public void appletResize(final int width, final int height) {
        this.minecraftApplet.resize(width, height);
    }

    @Override
    public void init() {
        this.minecraftApplet.init();
    }

    @Override
    public void start() {
        this.minecraftApplet.start();
    }

    @Override
    public void stop() {
        this.minecraftApplet.stop();
    }

    @Override
    public void destroy() {
        this.minecraftApplet.destroy();
    }

    @Override
    public void setVisible(final boolean visible) {
        super.setVisible(visible);
        this.minecraftApplet.setVisible(visible);
    }

}
