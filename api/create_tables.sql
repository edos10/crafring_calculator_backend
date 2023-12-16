CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
    -- Дополнительные характеристики предмета --
);

CREATE TABLE IF NOT EXISTS factories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    production_rate DOUBLE PRECISION NOT NULL
);

CREATE TABLE IF NOT EXISTS belts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
    -- Дополнительные характеристики ленты --
);

CREATE TABLE IF NOT EXISTS recipes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    item_id INT NOT NULL,
    factory_id INT NOT NULL,
    production_rate_per_factory DOUBLE PRECISION NOT NULL,
    FOREIGN KEY (item_id) REFERENCES items (id),
    FOREIGN KEY (factory_id) REFERENCES factories (id)
);


CREATE TABLE IF NOT EXISTS recipes_input (
       id SERIAL PRIMARY KEY,
       item_id INT NOT NULL,
       recipe_id INT NOT NULL,
       quantity INT NOT NULL,
       FOREIGN KEY (item_id) REFERENCES items (id),
       FOREIGN KEY (recipe_id) REFERENCES recipes (id)
);

CREATE TABLE IF NOT EXISTS recipe_belts (
    id SERIAL PRIMARY KEY,
    recipe_id INT NOT NULL,
    belt_id INT NOT NULL,
    quantity INT NOT NULL,
    FOREIGN KEY (recipe_id) REFERENCES recipes (id),
    FOREIGN KEY (belt_id) REFERENCES belts (id)
);

CREATE TABLE IF NOT EXISTS recipes_ierarchy
(
    id integer NOT NULL,
    child_id integer NOT NULL DEFAULT 0
)