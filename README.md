# battime
`battime` is a program that checks the battery usage time of your computer.

## Table of Contents
- [Features](#features)
- [How to use](#how-to-use)
- [Installation](#installation)
- [Updating](#updating)
- [License](#license)

## Feature
- Estimated Full Charge Time: See how much time is remaining until your battery is fully charged.
- Estimated Remaining Battery Time: Get an estimate of how much usage time you have left on your current battery charge.
- Battery Information: View detailed information about your device's battery.

## How to use
### Check battery status
```shell
battime
```
- When charging, Battery Charge Time is displayed in red.
- When not charging, Remaining Battery Time is displayed in cyan.

### View detailed battery information
```shell
battime --info / -i
```
- Your battery information (state, capacity, voltage, charge rate) is displayed.

## Installation
### On macOS
You can install `battime` with Homebrew:
```shell
brew tap chebread/battime

brew install battime
```

### For other OS
1. Visit [the GitHub Releases page](https://github.com/chebread/battime/releases) for `battime`.
2. Download the appropriate file for your operating system and architecture.
3. Unachive the downloaded file.
4. Execute the battime executable file.
5. For easier access, consider adding battime executable file to your system's PATH environment variable.

## Updating
### On macOS
If you installed `battime` using Homebrew, you can easily upgrade to the latest version when it's released.

```shell
brew upgrade battime
```

### For Other OS
For other OS, you will need to download the new version from [the GitHub Releases page](https://github.com/chebread/battime/releases) for `battime`.

Download the latest release for your system and replace your old executable file with the new one.

## License
MIT LICENSE &copy; 2025 Cha Haneum
