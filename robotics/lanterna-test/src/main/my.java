package main;
import java.io.IOException;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.concurrent.atomic.AtomicBoolean;

import com.googlecode.lanterna.TerminalPosition;
import com.googlecode.lanterna.TerminalSize;
import com.googlecode.lanterna.gui2.ActionListBox;
import com.googlecode.lanterna.gui2.BasicWindow;
import com.googlecode.lanterna.gui2.Borders;
import com.googlecode.lanterna.gui2.Button;
import com.googlecode.lanterna.gui2.Direction;
import com.googlecode.lanterna.gui2.LinearLayout;
import com.googlecode.lanterna.gui2.MultiWindowTextGUI;
import com.googlecode.lanterna.gui2.Panel;
import com.googlecode.lanterna.gui2.TextGUI;
import com.googlecode.lanterna.gui2.Window;
import com.googlecode.lanterna.gui2.WindowBasedTextGUI;
import com.googlecode.lanterna.gui2.WindowListener;
import com.googlecode.lanterna.gui2.dialogs.MessageDialog;
import com.googlecode.lanterna.gui2.dialogs.MessageDialogButton;
import com.googlecode.lanterna.input.KeyStroke;
import com.googlecode.lanterna.screen.Screen;
import com.googlecode.lanterna.screen.TerminalScreen;
import com.googlecode.lanterna.terminal.DefaultTerminalFactory;
import com.googlecode.lanterna.terminal.Terminal;


public class my {

    private static String showMenu(final WindowBasedTextGUI userInterface, String promptString, List<String> items) {

        // Notice that result is a final array containing a single String rather
        // than being a single, non-final String by itself.  That's because
        // Java's Runnables are not true closures, and the language only allows
        // them to access variables that are final or effectively final.  That
        // essentially prevents them from being able to modify variables in an
        // enclosing scope, unlike C#, C++, or JavaScript.
        //
        // But there's a loophole: in Java, you can make an array final and
        // still modify the value of items within the array.  It's a diry trick,
        // but I make no apologies.
        //
        // I learned the trick from https://stackoverflow.com/a/4732586.
        final String[] result = { items.get(0) };

        // Create a window that holds everything.
        final BasicWindow dialogWindow = new BasicWindow("Dialog Window");
        dialogWindow.setHints(Arrays.asList(Window.Hint.CENTERED));

        // The window must have one component, and it's going to be this one.
        Panel contentPanel = new Panel();
        dialogWindow.setComponent(contentPanel);

        // Our content consists in turn of two panels: one for the menu items
        // and another for the "Ok" and "Cancel" buttons.
        Panel mainPanel = new Panel();
        Panel buttonPanel = new Panel();
        contentPanel.addComponent(mainPanel.withBorder(Borders.singleLine("Main Panel")));
        contentPanel.addComponent(buttonPanel);

        // Construct the main panel.
        TerminalSize listBoxTerminalSize = new TerminalSize(20, 10);
        final ActionListBox listBox = new ActionListBox(listBoxTerminalSize);
        for (final String itemString : items) {
            // Populate the listBox with items from our string list.
            listBox.addItem(itemString, new Runnable() {
                public void run() {
                    // Selecting this list item just modifies the result
                    // variable.
                    //
                    // Sorry, result "array."
                    result[0] = itemString;
                }
            });
        }
        mainPanel.addComponent(listBox);

        // Construct the button panel.
        buttonPanel.setLayoutManager(new LinearLayout(Direction.HORIZONTAL));

        // We'll be using these callbacks in multiple places.
        final Runnable closeWindowHandler = new Runnable() {
            public void run() {
                dialogWindow.close();
            }
        };

        Button okButton = new Button("Ok", closeWindowHandler);
        Button cancelButton = new Button("Cancel", closeWindowHandler);
        buttonPanel.addComponent(okButton);
        buttonPanel.addComponent(cancelButton);

        // Handle keyboard input for the dialog window.
        //
        // I wouldn't have been able to implement this without https://stackoverflow.com/a/39124044.
        //
        // My question is why this keyboard listening code is not part of the official documentation.
        class KeyboardEventListener implements WindowListener {
            public boolean onUnhandledKeyStroke(TextGUI arg0, KeyStroke arg1) {
                return false;
            }

            public void onInput(Window window, KeyStroke key, AtomicBoolean pressed) {
                switch (key.getKeyType()) {
                    case Character:
                        switch (key.getCharacter()) {
                            case 'o':
                                // User hit the "Ok" button's hotkey.
                                closeWindowHandler.run();
                                break;
                            case 'c':
                                // User hit the "Cancel" button's hotkey.
                                closeWindowHandler.run();
                                break;
                        }
                        break;
                    case Enter:
                        if (listBox.isFocused()) {
                            // Treat this like selecting an individual list item.
                            listBox.getSelectedItem().run();
                            closeWindowHandler.run();
                        }
                        break;
                    case Escape:
                        // Same as hitting the "Cancel" button.
                        closeWindowHandler.run();
                        break;
                }
            }
            public void onUnhandledInput(Window arg0, KeyStroke arg1, AtomicBoolean arg2) { }
            public void onMoved(Window arg0, TerminalPosition arg1, TerminalPosition arg2) { }
            public void onResized(Window arg0, TerminalSize arg1, TerminalSize arg2) { }
        }
        dialogWindow.addWindowListener(new KeyboardEventListener());

        // Allow keyboard interaction with the dialog.
        listBox.takeFocus();

        // Place the window on top of the GUI's window stack.
        userInterface.addWindowAndWait(dialogWindow);

        // Label promptLabel = new Label(promptString);
        return result[0];
    }

    public static void main(String[] args) throws IOException {

        Terminal terminal = new DefaultTerminalFactory().createTerminal();
        Screen screen = new TerminalScreen(terminal);

        try {
            // Setup.
            screen.startScreen();

            // Create the initial window.
            WindowBasedTextGUI userInterface = new MultiWindowTextGUI(screen);
            List<String> items = new ArrayList<String>();
            items.add("One");
            items.add("Two");
            items.add("Three");
            String result = showMenu(userInterface, "Welcome to VRION", items);

            // The main window closed.
            //
            // Let's end by showing the user what they selected.
            MessageDialog.showMessageDialog(userInterface,
                    "MessageDialog",
                    String.format("You selected \"%s\"", result),
                    MessageDialogButton.OK);

        } catch (Exception e) {

            e.printStackTrace();

        } finally {

            // We're done!  Clean up the screen.
            try {
                if (screen != null) {
                    screen.stopScreen();
                }
                if (terminal != null) {
                    terminal.clearScreen();
                    terminal.resetColorAndSGR();
                }
            } catch (IOException e) {
                e.printStackTrace();
            }
        }

     /*outer:
     while(true) {
         KeyStroke keyStroke = screen.readInput();
     switch(keyStroke.getKeyType()){
     case Escape:
    	 break outer;
     }
     }*/


    }
}

