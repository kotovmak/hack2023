CREATE TABLE
    IF NOT EXISTS `z_app_tokens` (
        `ID` int(18) UNSIGNED NOT NULL AUTO_INCREMENT,
        `UF_USER_ID` int(18) DEFAULT NULL,
        `UF_APP_TOKEN` varchar(300) COLLATE utf8_unicode_ci DEFAULT NULL,
        PRIMARY KEY (`ID`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;