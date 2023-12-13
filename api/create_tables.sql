CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
    -- Дополнительные характеристики предмета --
);

CREATE TABLE factories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    production_rate DOUBLE PRECISION NOT NULL
);

CREATE TABLE belts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
    -- Дополнительные характеристики ленты --
);

CREATE TABLE recipes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    item_id INT NOT NULL,
    factory_id INT NOT NULL,
    production_rate_per_factory DOUBLE PRECISION NOT NULL,
    FOREIGN KEY (item_id) REFERENCES items (id),
    FOREIGN KEY (factory_id) REFERENCES factories (id)
);

CREATE TABLE recipe_belts (
    id SERIAL PRIMARY KEY,
    recipe_id INT NOT NULL,
    belt_id INT NOT NULL,
    quantity INT NOT NULL,
    FOREIGN KEY (recipe_id) REFERENCES recipes (id),
    FOREIGN KEY (belt_id) REFERENCES belts (id)
);
