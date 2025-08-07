# battime

[[한국어](README.kr.md)]

`battime` is a program that checks the battery usage time of your computer.

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
1. Visit [the GitHub Releases page](https://github.com/chebread/battime/releases) for `battime`.
2. Download the appropriate file for your operating system and architecture.
3. Unachive the downloaded file.
4. Execute the battime executable file.
5. For easier access, consider adding battime executable file to your system's PATH environment variable.

## License
MIT LICENSE &copy; 2025 Cha Haneum
