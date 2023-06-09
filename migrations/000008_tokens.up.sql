CREATE TABLE
    IF NOT EXISTS `z_app_tokens` (
        `ID` int(18) UNSIGNED NOT NULL AUTO_INCREMENT,
        `UF_USER_ID` int(18) DEFAULT NULL,
        `UF_APP_TOKEN` varchar(300) COLLATE utf8_unicode_ci DEFAULT NULL,
        PRIMARY KEY (`ID`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;

ALTER TABLE `z_prav_acts` ADD FULLTEXT (`UF_NAME`);

ALTER TABLE `z_nadzor_organs ` ADD FULLTEXT (`UF_NAME`);

ALTER TABLE `z_faq ` ADD FULLTEXT (`UF_QUESTION`,`UF_ANSWER`);

ALTER TABLE
    `sitemanager`.`z_notifications`
ADD
    COLUMN `UF_CONSULTATION_ID` int(18) NOT NULL;