#!/usr/bin/env python3
#
import re,time,sys,os


LOGFILE = 'sre_test_access.log.xz'
COUNTDOWN_NUM = 10
REGEX = r'\"(.*?)\"|\[(.*?)\]|(\S+)'
MSG = "Usage: " + "\t" + sys.argv[0] + "\t" + "< log filename >"
class LogMsgParser:


    def __init__(self, regex, rawlog):
        self.blist = []
        self.rawlog = rawlog
        self._regex_pattern = regex

    def logToArray(self):
        try:
            with open(self.rawlog, 'r', encoding='utf8') as f1:
                for i in f1.readlines():
                    if i == "\n":
                        continue
                    else:
                        self.blist.append(list(map(''.join, re.findall(self._regex_pattern, i))))
        except IOError as e:
            print("Log file not found or wrong format -> ", e.filename)

        return self.blist


class LogMsgHandler():

    def countHTTPReq(self):
        HttpReqs = []
        instance = LogMsgParser(REGEX, LOGFILE)
        b1 = instance.logToArray()
        for i in b1:
            try:
                HttpReqs.append(i)
            except IndexError as e:
                continue

        return HttpReqs.__len__()

    def _convertRequestTimeToEpoc(self, human_time):

        self.human_time = human_time
        epocTime = int(time.mktime(time.strptime(self.human_time, '%Y-%m-%d %H:%M:%S')))
        return epocTime


    def _convertLogTimeToEpoc(self, log_time):

        sliced_time = log_time.split()[0]
        epocTime = int(time.mktime(time.strptime(sliced_time, '%d/%b/%Y:%H:%M:%S')))

        return epocTime


    def _countIPAddr(self):

        ip_list, ip_count, total_list = [],[],[]
        ipObj = {}
        counter = 0
        g1 = ( i for i in self.timeFramedLog())
        for j in g1:
            ip_list.append(j[0])
        for k in ip_list:
            if k in ip_count:
                counter += 1
            else:
                ip_count.append(k)
            ipObj[k] = counter
        total_list.append(ipObj)
        del ip_count

        return total_list[0]    # return Dictionary object

    @property
    def mostHostFrom(self):
        x = self._countIPAddr()
        sorted_by_value = sorted(x.items(), key=lambda kv: kv[1], reverse=True)

        return sorted_by_value[:COUNTDOWN_NUM]     #return top 10 ip addr in reversed order.



    def timeFramedLog(self):

        timeframed_log = []
        instance = LogMsgParser(REGEX,LOGFILE)
        b1 = instance.logToArray() # First line is space
        for i in b1:
            if int(int(self._convertLogTimeToEpoc(i[3])) > self._convertRequestTimeToEpoc(start_t)) and int(self._convertLogTimeToEpoc(i[3]) < int(self._convertRequestTimeToEpoc(end_t))):
                timeframed_log.append(i)

        return timeframed_log


    def findOriginCountry(self, ipaddr):

        import geoip2.database
        reader = geoip2.database.Reader('GeoLite2-Country.mmdb')
        response = reader.country(ipaddr)

        return response.country.names['en']
start_t = os.getenv("STARTAT", default='2017-06-10 00:00:00') 
end_t = os.getenv("ENDAT", default='2017-06-19 23:59:59')
# Set STARTAT and ENDAT environment variable before running code.
# eg: export STARTAT='2017-01-01 00:00:00' 
#      Default time : Start at 2017-06-10 00:00:00 , End at 2017-06-19 23:59:59
                    

def main():
    #
    if len(sys.argv) <= 1:
        print("\n" + MSG + "\n")
        exit(1)
    LOGFILE = sys.argv[1]
    instance = LogMsgHandler()
    print("Totol number of HTTP request:")
    print(instance.countHTTPReq())

    result = instance.mostHostFrom
    print("Top 10 remote hosts accessed from " + start_t + " to " + end_t )
    print("IP Address\t\t", "Count\t\t", "Country" )
    for key in result:
        print(key[0] + '\t\t', key[1], '\t\t', instance.findOriginCountry(key[0]))


if __name__ == '__main__':
    main()
