// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package net.minecraft;

import java.applet.Applet;
import java.applet.AppletStub;
import java.awt.BorderLayout;
import java.net.MalformedURLException;
import java.net.URL;

public class Launcher extends Applet implements AppletStub {

    private static final String BASE = "http://www.minecraft.net/game/";

    private final Applet minecraftApplet;

    public Launcher(final Applet minecraftApplet) {
        this.minecraftApplet = minecraftApplet;
        this.minecraftApplet.setStub(this);

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

        try {
            return super.getParameter(name);
        }
        catch (final Exception ignored) {
            return null;
        }
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
