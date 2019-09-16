#!/bin/bash

MYSQL="mysql -uroot -h 127.0.0.1 --skip-column-names bte"

test1=`echo "select 1" | $MYSQL`

if [ $? -ne 0 ]; then
   echo "failed to connect to mysql"
   exit 1
fi

mtable="migration_log"

echo "CREATE TABLE IF NOT EXISTS $mtable (id varchar(20) not null primary key, action varchar(80) default '', run_at timestamp not null default CURRENT_TIMESTAMP) engine=InnoDB" | $MYSQL

if [ $? -ne 0 ]; then
   echo "failed to create $mtable"
   exit 1
fi

curdir=`pwd`
mpath=`dirname $0`
cd $mpath
migrations=`ls -1 -A *.sql`

for m in ${migrations[@]}; do
    id=${m:0:15}
    action=${m:16:-4}

    e=`echo "SELECT id FROM $mtable WHERE id='$id'" | $MYSQL`
    if [ $? -ne 0 ];then
       echo "failed to access $mtable"
       exit 1
    fi

    if [ "$e" "=" "" ]; then
       printf "running $m... "
       START=$(date +%s)
       cat $m | $MYSQL
       if [ $? -eq 0 ]; then
          END=$(date +%s)
          DIFF=$(( $END - $START ))
          echo "done in $DIFF seconds"
          echo "INSERT INTO $mtable set id='$id', action='$action'" | $MYSQL
       fi
    elif [ "$e" "!=" "$id" ]; then
       echo "skip $m"
    fi
done
cd $curdir

echo
echo "DONE"
