#!/bin/bash

function loading {
  sleep 1
  printf "."
  sleep 1
  printf "."
  sleep 1
  printf "."
  sleep 1
}

echo "Hi there!
I see you want to re/create pharmacodb from scratch ...
Well aren't you brave, attempting what many have tried, and failed ...
But fret not, for this script will make your life a lot simpler :D
Easy peasy right, so just follow along on this crazy ride!"
echo
printf "
                                                      \  /
                __                                     \/
   _   ---===##===---_________________________--------------  _
  [ ~~~=================###=###=###=###=###=================~~ ]
  /  ||  | |~\  ;;;;     PKP    ;;;  ET22-689  ;;;;  /~| |  ||  \\
 /___||__| |  \ ;;;;            [_]            ;;;; /  | |__||___\\
 [\        |__| ;;;;  ;;;; ;;;; ;;; ;;;; ;;;;  ;;;; |__|        /]
(=|    ____[-]_______________________________________[-]____Kraq|=)
/  /___/|#(__)=o########o=(__)#||___|#(__)=o#########o=(__)#|\___\\
_________-=\__/=--=\__/=--=\__/=-_____-=\__/=--=\__/=--=\__/=-______

"
echo
echo

echo "You are not under oath, but do answer each question as honestly as you can!"
echo "Oh, and if you fail to comply with any instructions,
script will get angry at you and exit rudely ..."
echo

# printf "starting "
# loading
# echo
# echo

echo "Do you have a mysql user and password already created?"
echo "  (If you enter no, you will be prompted to create a user.)"
echo "  Press y for Yes, n for No."
printf "(y/n): "
read ans
echo

if [[ ans -eq "y" ]]; then
  printf ""



  exit 0
elif [[ ans -eq "n" ]]; then
  echo "no"
  exit 0
elif [[ ans -eq "h" ]]; then
  echo "help"
  exit 0
else
  echo "You failed to comply with simple instructions :("
  printf "exiting non-gracefully "
  loading
  echo
  exit 1
fi

printf "Enter user name: "
read name
printf "Enter password: "
stty -echo
read password
stty echo
printf "\n"
echo "username is: $name, password is: $password"
