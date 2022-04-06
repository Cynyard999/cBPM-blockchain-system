package com.cbpm.backend.dao;

import com.cbpm.backend.vo.UserVo;
import org.springframework.data.repository.CrudRepository;

import java.util.List;

public interface UserRepository extends CrudRepository<UserVo, Integer> {

    List<UserVo> findByOrgType(String orgType);

    UserVo findByEmail(String email);

}
