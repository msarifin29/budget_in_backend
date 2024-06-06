DROP TABLE IF EXISTS `users`; 
DROP TABLE IF EXISTS `expenses`; 
DROP TABLE IF EXISTS `incomes`; 
DROP TABLE IF EXISTS `accounts`; 
ALTER TABLE IF EXISTS "t_category_expenses" DROP CONSTRAINT IF EXISTS "fk_t_category_expenses_expenses";
ALTER TABLE IF EXISTS "t_category_incomes" DROP CONSTRAINT IF EXISTS "fk_t_category_incomes_incomes";
ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "fk_accounts_users";
