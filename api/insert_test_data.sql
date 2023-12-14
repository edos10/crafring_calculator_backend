INSERT INTO recipes (
	id, name, item_id, factory_id, production_rate_per_factory)
	VALUES (1, "Water Production", 1, 1, 72000);

INSERT INTO recipe_belts (
	id, recipe_id, belt_id, quantity)
	VALUES (1, 1, 1, 17);


INSERT INTO items (
	id, name)
	VALUES (1, "water");

INSERT INTO factories(
	id, name, production_rate)
	VALUES (1, "offshore pump", 72000);

INSERT INTO belts(
	id, name)
	VALUES (1, "pipes");
    