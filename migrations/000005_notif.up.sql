CREATE TABLE
    IF NOT EXISTS `z_notifications` (
        `ID` int(18) UNSIGNED NOT NULL AUTO_INCREMENT,
        `UF_DATE` datetime NOT NULL,
        `UF_TEXT` text DEFAULT NULL,
        `UF_USER_ID` int(18) NOT NULL,
        PRIMARY KEY (`ID`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;

INSERT INTO
    `hack2023`.`z_notifications` (
        `UF_DATE`,
        `UF_TEXT`,
        `UF_USER_ID`
    )
VALUES (
        '2023-05-26 14:24:54',
        'some text',
        1
    ), (
        '2023-05-27 14:24:54 ',
        'some text',
        1
    );