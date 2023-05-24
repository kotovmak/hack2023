CREATE TABLE
    IF NOT EXISTS `z_slots` (
        `ID` int(18) UNSIGNED NOT NULL AUTO_INCREMENT,
        `UF_DATE` date NOT NULL,
        `UF_TIME` varchar(50) NOT NULL,
        PRIMARY KEY (`ID`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;

CREATE TABLE
    IF NOT EXISTS `z_consultations` (
        `ID` int(18) UNSIGNED NOT NULL AUTO_INCREMENT,
        `UF_DATE` date NOT NULL,
        `UF_TIME` varchar(50) NOT NULL,
        `UF_QUESTION` text DEFAULT NULL,
        `UF_NADZOR_ORGAN_ID` int(18) NOT NULL,
        `UF_CONTROL_TYPE_ID` int(18) NOT NULL,
        `UF_CONSULT_TOPIC_ID` int(18) NOT NULL,
        `UF_USER_ID` int(18) NOT NULL,
        `UF_IS_NEED_LATTER` BOOLEAN DEFAULT NULL,
        `UF_IS_CONFIRMED` BOOLEAN DEFAULT NULL,
        PRIMARY KEY (`ID`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;