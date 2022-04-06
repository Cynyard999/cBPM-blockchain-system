package com.cbpm.backend.vo;

import javax.persistence.*;

/**
 * @author Polaris
 * @version 1.0
 * @description: 实际user类
 * @date 2022/4/5 16:35
 */

@Table(name = "user")
@Entity
public class UserVo {

    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)
    private Integer id;
    private String name;
    private String password;
    private String orgType;
    private String email;

    public UserVo() {

    }


    @Override
    public String toString() {
        return "{" +
                " id='" + getId() + "'" +
                ", userName='" + getName() + "'" +
                ", email='" + getEmail() + "'" +
                ", password='" + getPassword() + "'" +
                ", orgtype='" + getOrgType() +
                "}";
    }

    public UserVo(String name, String password, String orgType, String email) {
        this.name = name;
        this.password = password;
        this.orgType = orgType;
        this.email = email;
    }

    public UserVo(Integer id, String name, String password, String orgType, String email) {
        this.id = id;
        this.name = name;
        this.password = password;
        this.orgType = orgType;
        this.email = email;
    }

    public Integer getId() {
        return id;
    }

    public void setId(Integer id) {
        this.id = id;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getPassword() {
        return password;
    }

    public void setPassword(String password) {
        this.password = password;
    }

    public String getEmail() {
        return email;
    }

    public void setEmail(String email) {
        this.email = email;
    }

    public String getOrgType() {
        return orgType;
    }

    public void setOrgType(String orgType) {
        this.orgType = orgType;
    }
}
