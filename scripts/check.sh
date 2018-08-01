#!/bin/bash

if [[ $# == 0 || $# -gt 3 ]]; then
  echo "ERROR: Incorrect number of parameters"
  echo "Usage: $0 [-v] dir1 [dir2]"
  exit 1
fi

verbose=
dir1=
dir2=
while (( "$#" )); do 
  case $1 in
    "-v") 
      verbose=true
      ;;
    *)
      if [[ -z $dir1 ]]; then
        dir1=$1
      else 
        dir2=$1
      fi
      ;;
  esac
  shift
done

[[ -z $dir2 ]] && dir2="/Volumes/SMC/Shared iTunes/iTunes Media/Movies"

if [[ -z $dir1 ]]; then
  echo "ERROR: Missing directory 1"
  exit 1
fi

cd "$dir1" # "/Volumes/SMC/_To Organize/Movies/To_Delete_2"
for f in *; do 
  fname=${f##*/}; 
  title=${fname%.*};
  title=${title% (1080p HD)} 
  title=${title% (HD)} 
  title=$(echo $title | sed 's/\.$/_/')
  otherf="$dir2/$title/$fname";

  mtitle="Movie '$title'"

  if [[ ! -f "$otherf" ]]; then 
    echo "[ERROR ] $mtitle does not exist in iTunes ($otherf)"
    continue; 
  fi
  
  eval $(stat -s "$f"); 
  s1=$st_size;
  st_size=0
  
  eval $(stat -s "$otherf"); 
  s2=$st_size; 
  
  if [[ $s1 == $s2 ]]; then
    [[ $verbose == true ]] && echo "[  OK  ] $mtitle are the same size ($s1 == $s2)";
  else 
    echo "[ERROR ] $mtitle are different size ($s1 == $s2)"
  fi
done