#!/bin/sh

Tmp='./mysl'
MysqlSockList=`/bin/ps aux | /bin/grep mysql.sock | /bin/grep -v grep | /bin/awk '{print $(NF-1)}' | /bin/awk '{split($0,a,"=");print a[2]}' > $Tmp`
/bin/sed '/^\s*$/d' $Tmp > /dev/null
Yesterday=`/bin/date -d yesterday +"%Y%m%d"`

/bin/cat $Tmp | while read sock
do
        #echo $sock
        Slowlog=$(mysql -S $sock -Bse "show GLOBAL VARIABLES like 'slow_query_log_file';" | awk '{print $2}')

        if [ -f $Slowlog ];then
                \/bin/cp -rf $Slowlog $Slowlog.$Yesterday
                /bin/echo > $Slowlog
        fi
done

/bin/rm -rf $Tmp
