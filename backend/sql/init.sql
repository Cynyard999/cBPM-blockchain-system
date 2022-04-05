DROP TABLE IF EXISTS user_vo ;
DROP TABLE IF EXISTS user ;
CREATE TABLE `user` (
                      `id` int(11) NOT NULL AUTO_INCREMENT,
                      `name` char(50) NOT NULL,
                      `password` char(50) NOT NULL,
                      `org_type` char(50) NOT NULL,
                      `email` char(50) NOT NULL,
                      PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
INSERT INTO user (name, password, org_type, email) VALUES ('carrier-admin','adminPW','carrier','');
INSERT INTO user (name, password, org_type, email) VALUES ('supplier-admin','adminPW','supplier','');
INSERT INTO user (name, password, org_type, email) VALUES ('middleman-admin','adminPW','middleman','');
INSERT INTO user (name, password, org_type, email) VALUES ('manufacturer-admin','adminPW','manufacturer','');