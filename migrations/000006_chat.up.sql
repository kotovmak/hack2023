CREATE TABLE
    IF NOT EXISTS `z_messages` (
        `ID` int(18) UNSIGNED NOT NULL AUTO_INCREMENT,
        `UF_DATE` datetime NOT NULL,
        `UF_TEXT` text DEFAULT NULL,
        `UF_SEND_BY_ID` int(18) NOT NULL,
        `UF_USER_ID` int(18) NOT NULL,
        PRIMARY KEY (`ID`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;

CREATE TABLE
    IF NOT EXISTS `z_chat_buttons` (
        `ID` int(18) UNSIGNED NOT NULL AUTO_INCREMENT,
        `UF_TEXT` text DEFAULT NULL,
        `UF_LINK` varchar(200) NOT NULL,
        PRIMARY KEY (`ID`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;

INSERT INTO
    `z_chat_buttons` (`UF_TEXT`, `UF_LINK`)
VALUES (
        'Нормативные акты',
        'https://knd.mos.ru/npa'
    );

INSERT INTO
    `z_chat_buttons` (`UF_TEXT`, `UF_LINK`)
VALUES (
        'Обязательные требования',
        'https://knd.mos.ru/requirements/public'
    );

INSERT INTO
    `z_chat_buttons` (`UF_TEXT`, `UF_LINK`)
VALUES (
        'Информационные материалы',
        'https://knd.mos.ru/infomaterials'
    );