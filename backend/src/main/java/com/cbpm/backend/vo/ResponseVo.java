package com.cbpm.backend.vo;

/**
 * @author Polaris
 * @version 1.0
 * @description: 封装所有response
 * @date 2022/4/1 21:02
 */
public class ResponseVo {
    boolean successful;

    String message;

    Object content;

    public ResponseVo(boolean successful) {
        this.successful = successful;
    }

    public ResponseVo(String message) {
        this.message = message;
    }

    public ResponseVo(Object object) {
        this.content = object;
    }

    public boolean isSuccessful() {
        return successful;
    }

    public  static ResponseVo buildSuccess(){
        ResponseVo responseVo=new ResponseVo("success");
        responseVo.setSuccessful(true);
        return responseVo;
    }

    public static ResponseVo buildSuccess(String message){
        if(message.length()==0){
            message="success";
        }
        ResponseVo responseVo=new ResponseVo(true);
        System.out.println(message);
        responseVo.setMessage(message);

        return responseVo;
    }

    public static ResponseVo buildSuccess(Object content){
        ResponseVo responseVo=new ResponseVo(content);
        responseVo.setSuccessful(true);
        return responseVo;
    }

    public static ResponseVo buildFailure(String message){
        ResponseVo responseVo=new ResponseVo(message);
        System.out.println(message);
        responseVo.setSuccessful(false);
        return  responseVo;
    }

    public void setSuccessful(boolean successful) {
        this.successful = successful;
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }

    public Object getContent() {
        return content;
    }

    public void setContent(Object content) {
        this.content = content;
    }
}
