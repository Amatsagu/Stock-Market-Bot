# Reset just in case
if [ "$1" != "-reset"  ] || [ "$1" != "-r"  ]; then
    echo "Clearing old configuration..."
    systemctl stop index-bot.service
    rm "/lib/systemd/system/index-bot.service"
    systemctl daemon-reload
    journalctl --rotate
    journalctl --vacuum-time=1s
    chmod 777 "/etc/vironicer" -R
    rm /etc/vironicer/IndexBot/service/index-bot.o
fi

# Configure systemctl daemon
echo "Compiling..."
/usr/local/go/bin/go build -o /etc/vironicer/IndexBot/service/index-bot.o ./main.go
chmod 777 "/etc/vironicer" -R

echo "Registering new daemon..."
cp "/etc/vironicer/IndexBot/service/index-bot.service" "/lib/systemd/system/index-bot.service"
chmod 644 "/lib/systemd/system/index-bot.service"
systemctl daemon-reload
systemctl enable index-bot.service
systemctl restart index-bot
echo "Finished!"