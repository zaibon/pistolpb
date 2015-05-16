# pistolPB

little script that allow you to send some pushbullet notification when you receive a highlight on irc.

## Install
```go get github.com/zaibon/pistolpb```

## Usage
you need to enable the [fnotify](https://gist.github.com/matthutchinson/542141) script into irssi

then load the pistol with your api key and choose on which devices you want to shoot
```bash
pistolpb -k <you api key>
0 : samsung GT-S7275R
1 : Chrome
choose on which device send notification (ex: 0,1,3) : 
```

finally start shooting
```
pistolpb -f /path/to/fnotify/file