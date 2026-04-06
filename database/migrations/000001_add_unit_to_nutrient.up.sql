ALTER TABLE nutrient
    ADD COLUMN unit VARCHAR(20) NOT NULL DEFAULT 'g';

-- Set correct units for the 10 seeded nutrients
UPDATE nutrient SET unit = 'kcal' WHERE id = 1;  -- Energy
UPDATE nutrient SET unit = 'g'    WHERE id = 2;  -- Protein
UPDATE nutrient SET unit = 'g'    WHERE id = 3;  -- Fat
UPDATE nutrient SET unit = 'g'    WHERE id = 4;  -- Carbohydrate
UPDATE nutrient SET unit = 'g'    WHERE id = 5;  -- Fiber
UPDATE nutrient SET unit = 'mg'   WHERE id = 6;  -- Cholesterol
UPDATE nutrient SET unit = 'mcg'  WHERE id = 7;  -- Vitamin A
UPDATE nutrient SET unit = 'mcg'  WHERE id = 8;  -- Vitamin B12
UPDATE nutrient SET unit = 'mg'   WHERE id = 9;  -- Calcium
UPDATE nutrient SET unit = 'mg'   WHERE id = 10; -- Iron
