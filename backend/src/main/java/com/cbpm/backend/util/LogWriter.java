package com.cbpm.backend.util;

import java.sql.Timestamp;
import java.io.File;
import java.io.FileWriter;
import java.io.IOException;

import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Configuration;

/**
 * @author Polaris
 * @version 1.0
 * @description: queryLog
 * @date 2022/4/6 15:50
 */
@Slf4j
@Configuration

public class LogWriter {


    @Value("${backend.logPackagePath}")
    private String logPackagePath;


    public void writeLog(String logStr)  {
        try {
            System.out.println("logPath:"+this.logPackagePath);
            File file = new File(this.logPackagePath+"log.txt");
            if (!file.exists()) {
                file.createNewFile();
            }
            FileWriter fileWritter = new FileWriter(file, true);
            Timestamp timestamp = new Timestamp(System.currentTimeMillis());
            System.out.println(timestamp);
            fileWritter.write(timestamp+":"+"\n"+logStr+"\n");
            fileWritter.close();
        }catch (IOException e){
            System.out.println("write to log.txt wrong.");
        }

    }

}
