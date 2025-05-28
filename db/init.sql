DROP SCHEMA IF EXISTS rikonscraper CASCADE;

CREATE SCHEMA IF NOT EXISTS rikonscraper;

CREATE TABLE IF NOT EXISTS rikonscraper.products(
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    display text,
    url text,
    product_code text,
    CONSTRAINT product UNIQUE (product_code)
);

CREATE TABLE IF NOT EXISTS rikonscraper.parts(
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    parent UUID references rikonscraper.products(id),
    display text,
    url text,
    product_code text,
    CONSTRAINT product_parent UNIQUE (product_code, parent)
);

CREATE TABLE IF NOT EXISTS rikonscraper.alerts(
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    part_id UUID NOT NULL references rikonscraper.parts(id),
    email text,
    CONSTRAINT email_part UNIQUE (email, part_id)
);