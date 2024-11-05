# procezz
Watches processes on MacOS and can perform actions based on startup or shutdown of them. Below is how to set it up.

## Compile

```
go build -ldflags="-s -w" -o procezz .
```

## Running without system

```
make start
```

## Running on startup in MacOS

```
cat ~/Library/LaunchAgents/procezz.plist

<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>io.yourname.procezz</string>
    <key>ProgramArguments</key>
    <array>
      <string>/usr/local/bin/procezz</string>
      <string>/etc/procezz.conf</string>
    </array>
    <key>KeepAlive</key>
    <true/>
</dict>
</plist>
```

Now we need to check our new LaunchAgent. *! Do not start the agent ut=ntil you habe move the files or updated the script!*

```
launchctl list io.yourname.procezz
launchctl start ~/Library/LaunchAgents/procezz.plist
launchctl load ~/Library/LaunchAgents/procezz.plist
launchctl unload ~/Library/LaunchAgents/procezz.plist
launchctl reload ~/Library/LaunchAgents/procezz.plist
launchctl status io.yourname.procezz
```

## Move Files to Right Places or Update the Above LaunchAgent with your Directories

We need to move our tool to `/usr/local/bin`. Additionally we need to add the config file to `/etc/procezz.conf`

```
sudo cp tmp/procezz /usr/local/bin/
sudo cp procezz.conf /etc/
```

## Configure Process Watch

Make sure that your shell scripts are executable with `chmod +x some_shell_script_on_start`.

```
cat /etc/prosezz.conf

{
  "services": [
    {
      "name": "iTerm2",
      "on_start": "/Users/<your_user>/dev/bin/some_shell_script_on_start",
      "on_stop": "/Users/<your_user>/dev/bin/some_shell_script_on_stop"
    }
  ]
}

```
