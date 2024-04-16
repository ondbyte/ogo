CREATE TABLE categories (
  id INTEGER AUTO_INCREMENT PRIMARY KEY,
  name TEXT UNIQUE
);

CREATE TABLE tags (
    id INTEGER AUTO_INCREMENT PRIMARY KEY,
    name TEXT UNIQUE
);

CREATE TABLE pets (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  photo_url TEXT,
  status TEXT NOT NULL,
  category INTEGER,
  FOREIGN KEY (category) REFERENCES categories(id) 
);

CREATE TABLE pet_tags (
    pet_id INTEGER NOT NULL,
    tag_id INTEGER NOT NULL,
    FOREIGN KEY (pet_id) REFERENCES pets(id) ,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ,
    PRIMARY KEY (pet_id, tag_id)
);
