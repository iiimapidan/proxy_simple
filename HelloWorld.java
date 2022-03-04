import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.Calendar;
import java.util.Date;
import java.util.GregorianCalendar;

public class HelloWorld {

    public static boolean isBeforeCurDay(long t) throws ParseException {
        System.out.println(t);
        SimpleDateFormat dateFormat = new SimpleDateFormat("yyyy-MM-dd");

        Calendar calendar = new GregorianCalendar();
        calendar.setTimeInMillis(t * 1000);
        Date date = calendar.getTime();
        String strDate = dateFormat.format(date);
        Date date2 = dateFormat.parse(strDate);
        long t2 = date2.getTime();

        String curDate = dateFormat.format(new Date());
        Date curDate2 = dateFormat.parse(curDate);
        long curDate3 = curDate2.getTime();

        if (t2 < curDate3) {
            return true;
        }

        return false;
    }

    public static void main(String[] args) throws ParseException {
        // 2022.03.03 21:24
        long recordTime = 1646402242;
        System.out.println(isBeforeCurDay(recordTime));
    }
}