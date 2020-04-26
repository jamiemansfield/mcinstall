// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

import net.minecraftforge.installer.SimpleInstaller;
import net.minecraftforge.installer.actions.ActionCanceledException;
import net.minecraftforge.installer.actions.ClientInstall;
import net.minecraftforge.installer.actions.ProgressCallback;
import net.minecraftforge.installer.json.Install;
import net.minecraftforge.installer.json.Util;

import java.io.File;

public final class ModernForgeClientTool {

    public static void main(final String[] args) {
        if (args.length < 1) {
            System.err.println("must provide launcher directory!");
            System.exit(-1);
        }
        final File installDir = new File(args[0]);

        SimpleInstaller.headless = true;
        final Install profile = Util.loadInstallProfile();
        final ClientInstall install = new ClientInstall(profile, ProgressCallback.TO_STD_OUT);
        try {
            if (!install.run(installDir, op -> true)) {
                System.err.println("failed to install!");
                System.exit(-1);
            }
        }
        catch (final ActionCanceledException ex) {
            ex.printStackTrace();
            System.exit(-1);
        }
    }

    private ModernForgeClientTool() {
    }

}
