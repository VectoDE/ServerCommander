# ServerCommander Themes Guide

## Overview

ServerCommander allows users to customize the look and feel of the command-line interface by modifying the theme. This guide will explain how to apply and create custom themes for your environment.

## Default Theme

ServerCommander comes with a default theme that is designed to be simple, clean, and easy to read. It uses a combination of basic terminal colors to highlight important information.

## Applying a Theme

To change the theme, you need to modify the configuration file or use the `--theme` option when running the command.

### 1. Using the `--theme` Flag

You can temporarily apply a theme for a single session by passing the `--theme` flag followed by the theme name. For example:

```bash
server-commander --theme dark
```

This command will apply the ```dark``` theme for the current session. The available theme names are ```dark```, ```light```, and ```custom```.

### 2. Changing the Default Theme

To set a default theme for all sessions, you need to modify the configuration file (```config.yaml```) and specify the desired theme.

```yaml
theme: dark
```

Save the configuration, and the theme will be applied automatically for all future sessions.

## Creating a Custom Theme

If you want to create a custom theme, you can do so by defining your colors in the configuration file. Each theme is a set of key-value pairs that represent different elements of the interface, such as text colors, background colors, and highlights.

### Example Custom Theme

```yaml
theme: custom
colors:
  background: "#1c1c1c"
  text: "#f8f8f2"
  primary: "#ff79c6"
  secondary: "#8be9fd"
  highlight: "#50fa7b"
```

This theme changes the background to a dark shade, the text to a light color, and defines custom colors for primary, secondary, and highlighted text.

### Key Color Elements

Here are the main elements you can customize in your theme:

- **background**: The background color of the terminal.
- **text**: The color of the regular text.
- **primary**: The color for primary information, such as server names or commands.
- **secondary**: The color for secondary information, such as paths or commands output.
- **highlight**: The color used to highlight important text or active commands.

Once you have defined your custom colors, save the ```config.yaml``` file and restart ServerCommander to apply the new theme.

### Predefined Themes

ServerCommander includes a few predefined themes:

- **dark**: A dark theme with a modern, minimal design.
- **light**: A light theme for users who prefer a bright interface.
- **custom**: A theme that you can fully customize with your own colors.

## Conclusion

Themes allow you to personalize your experience with ServerCommander, making it easier on the eyes and more enjoyable to use. You can always switch between themes or create your own custom styles to match your preferences.

For more information on configuring ServerCommander, refer to the [Configuration Guide](CONFIGURATION.md).
