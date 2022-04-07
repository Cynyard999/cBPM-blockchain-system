package com.cbpm.backend.util;

import java.sql.Timestamp;
import java.io.File;
import java.io.FileWriter;
import java.io.IOException;

import java.text.SimpleDateFormat;
import java.util.Date;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Configuration;

/**
 * @author Polaris, cynyard
 * @version 1.1
 * @date 2022/4/7 9:50
 */
@Slf4j
@Configuration

public class LogWriter {


    @Value("${backend.logPackagePath}")
    private String logPackagePath;


    public void writeLog(String logStr) {
        try {
            Long timeStamp = System.currentTimeMillis();
            String date = new SimpleDateFormat("yyyy-MM-dd")
                    .format(new Date(Long.parseLong(String.valueOf(timeStamp))));
            File file = new File(this.logPackagePath + date + ".log");
            if (!file.exists()) {
                file.createNewFile();
            }
            FileWriter fileWriter = new FileWriter(file, true);
            fileWriter.write(new Timestamp(System.currentTimeMillis()).toString() + ": " + logStr
                    + "\n");
            fileWriter.close();
        } catch (IOException e) {
            System.out.println("writing to log failed.");
        }

    }

}
