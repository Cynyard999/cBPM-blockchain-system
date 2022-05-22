package com.cbpm.backend.vo;

/**
 * @author Polaris, Cynyard
 * @version 1.1
 * @description: 封装所有response
 * @date 2022/4/1 21:02
 */
public class ResponseVo {

    boolean success;

    String message;

    Object result;

    int status;

    private ResponseVo(boolean success) {
        this.success = success;
    }

    private ResponseVo(String message) {
        this.message = message;
    }

    private ResponseVo(Object object) {
        this.result = object;
    }

    public ResponseVo(int status) {
        this.status = status;
    }

    private ResponseVo(String message, Object result) {
        this.message = message;
        this.result = result;
    }

    public boolean isSuccess() {
        return success;
    }

    public static ResponseVo buildSuccess(String message, Object content) {
        if (message == null || message.length() == 0) {
            message = "success";
        }
        ResponseVo responseVo = new ResponseVo(message, content);
        responseVo.setSuccess(true);
        responseVo.setStatus(200);
        return responseVo;
    }

    public static ResponseVo buildFailure(String message, int status) {
        if (message == null || message.length() == 0) {
            message = "failure";
        }
        ResponseVo responseVo = new ResponseVo(message);
        responseVo.setSuccess(false);
        responseVo.setStatus(status);
        return responseVo;
    }

    public static ResponseVo buildFailure(String message, Object content, int status) {
        if (message == null || message.length() == 0) {
            message = "failure";
        }
        ResponseVo responseVo = new ResponseVo(message, content);
        responseVo.setSuccess(false);
        responseVo.setStatus(status);
        return responseVo;
    }

    public void setSuccess(boolean success) {
        this.success = success;
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }

    public Object getResult() {
        return result;
    }

    public void setResult(Object result) {
        this.result = result;
    }

    public int getStatus() {
        return status;
    }

    public void setStatus(int status) {
        this.status = status;
    }

    @Override
    public String toString() {
        String res = "";
        if (this.success) {
            res += "success,";
        } else {
            res += "false,";
        }
        return res + "message:" + this.getMessage();
    }
}
