#!/bin/sh

chown tor -R /var/lib/tor 
chmod 700 -R /var/lib/tor

/usr/bin/tor -f /etc/tor/torrc --verify-config
if [ "$?" != "0" ]; then
  echo "INVALID TOR CONFIG!!"
  exit 1
fi

/usr/bin/tor -f /etc/tor/torrc
sleep 1

echo
echo "======== ACCESS THE INTERFACE AT THIS ADDRESS ========"
cat /var/lib/tor/venom/hostname
echo "============== NOW STARTING THE SERVER ==============="
echo

./server
