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

    public ResponseVo(boolean success) {
        this.success = success;
    }

    public ResponseVo(String message) {
        this.message = message;
    }

    public ResponseVo(Object object) {
        this.result = object;
    }

    public boolean isSuccess() {
        return success;
    }

    public static ResponseVo buildSuccess() {
        ResponseVo responseVo = new ResponseVo("success");
        responseVo.setSuccess(true);
        return responseVo;
    }

    public static ResponseVo buildSuccess(String message) {
        if (message.length() == 0) {
            message = "success";
        }
        ResponseVo responseVo = new ResponseVo(true);
        System.out.println(message);
        responseVo.setMessage(message);

        return responseVo;
    }

    public static ResponseVo buildSuccess(Object content) {
        ResponseVo responseVo = new ResponseVo(content);
        responseVo.setSuccess(true);
        return responseVo;
    }

    public static ResponseVo buildFailure(String message) {
        ResponseVo responseVo = new ResponseVo(message);
        System.out.println(message);
        responseVo.setSuccess(false);
        return responseVo;
    }

    public static ResponseVo buildFailure(Object content) {
        ResponseVo responseVo = new ResponseVo(false);
        responseVo.setResult(content);
        return responseVo;
    }

    public static ResponseVo buildFailure(String[] exceptions) {
        ResponseVo responseVo = new ResponseVo(false);
        StringBuilder str = new StringBuilder();
        for (int i = 1; i < exceptions.length; i++) {
            str.append(exceptions[i]);
        }
        responseVo.setMessage(str.toString());
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
}
