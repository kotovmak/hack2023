CREATE TABLE
    IF NOT EXISTS `z_api_users` (
        `ID` int(18) NOT NULL AUTO_INCREMENT,
        `UF_LOGIN` varchar(50) COLLATE utf8_unicode_ci NOT NULL,
        `UF_PASSWORD` varchar(100) COLLATE utf8_unicode_ci NOT NULL,
        `UF_ACTIVE` char(1) COLLATE utf8_unicode_ci NOT NULL DEFAULT 'Y',
        `UF_NAME` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
        `UF_EMAIL` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
        `UF_DATE_REGISTER` datetime NOT NULL,
        PRIMARY KEY (`ID`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;

CREATE TABLE
    IF NOT EXISTS `z_refresh_tokens` (
        `ID` int(18) UNSIGNED NOT NULL AUTO_INCREMENT,
        `UF_USER_ID` int(18) DEFAULT NULL,
        `UF_REFRESH_TOKEN` varchar(200) COLLATE utf8_unicode_ci DEFAULT NULL,
        `UF_UPDATED_AT` datetime NOT NULL,
        PRIMARY KEY (`ID`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;

CREATE TABLE
    IF NOT EXISTS `z_nadzor_organs` (
        `ID` int(18) UNSIGNED NOT NULL AUTO_INCREMENT,
        `UF_NAME` varchar(200) COLLATE utf8_unicode_ci NOT NULL,
        PRIMARY KEY (`ID`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;

CREATE TABLE
    IF NOT EXISTS `z_services` (
        `ID` int(18) UNSIGNED NOT NULL AUTO_INCREMENT,
        `UF_NAME` varchar(200) COLLATE utf8_unicode_ci NOT NULL,
        `UF_DESCRIPTION` varchar(500) COLLATE utf8_unicode_ci DEFAULT NULL,
        PRIMARY KEY (`ID`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;

CREATE TABLE
    IF NOT EXISTS `z_control_types` (
        `ID` int(18) UNSIGNED NOT NULL AUTO_INCREMENT,
        `UF_NAME` varchar(300) COLLATE utf8_unicode_ci NOT NULL,
        `UF_NADZOR_ORGAN_ID` int(18) NOT NULL,
        PRIMARY KEY (`ID`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;

CREATE TABLE
    IF NOT EXISTS `z_consult_topics` (
        `ID` int(18) UNSIGNED NOT NULL AUTO_INCREMENT,
        `UF_NAME` varchar(300) COLLATE utf8_unicode_ci NOT NULL,
        `UF_NADZOR_ORGAN_ID` int(18) NOT NULL,
        `UF_CONTROL_TYPE_ID` int(18) NOT NULL,
        PRIMARY KEY (`ID`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;

CREATE TABLE
    IF NOT EXISTS `z_prav_acts` (
        `ID` int(18) UNSIGNED NOT NULL AUTO_INCREMENT,
        `UF_NAME` varchar(200) COLLATE utf8_unicode_ci NOT NULL,
        `UF_NADZOR_ORGAN_ID` int(18) NOT NULL,
        `UF_CONTROL_TYPE_ID` int(18) NOT NULL,
        PRIMARY KEY (`ID`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;