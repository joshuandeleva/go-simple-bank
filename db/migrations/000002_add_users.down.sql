ALTER TABLE IF EXISTS "account" DROP CONSTRAINT IF EXISTS "accounts_owner_fkey";


ALTER TABLE IF EXISTS "account" DROP CONSTRAINT IF EXISTS "owner_currency_key";


DROP TABLE IF EXISTS "users";

