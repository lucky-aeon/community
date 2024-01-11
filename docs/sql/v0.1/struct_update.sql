-- change invite_code table: code field type to int(8)
-- dateTime: 2024-01-11 21:07
ALTER TABLE `invite_code`
	CHANGE `code` `code` INT(8) DEFAULT NULL ;