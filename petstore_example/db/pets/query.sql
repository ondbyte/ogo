-- name: AddCategory :execresult
insert into categories(name) values (?);

-- name: GetCategoryByName :one
select * from categories where categories.name = ?;

-- name: AddTag :execresult
INSERT INTO tags(name) values (?);

-- name: GetTagByName :one
select * from tags where tags.name = ?;

-- name: AddPet :execresult
insert into pets(name,photo_url,status,category) values(?,?,?,?);

-- name: AddPetTag :execresult
insert into pet_tags(pet_id,tag_id) values(?,?)