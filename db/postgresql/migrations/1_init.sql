CREATE TABLE users
(
                       id int GENERATED ALWAYS AS IDENTITY NOT NULL,
                       login text NOT NULL,
                       password text NOT NULL,
                       CONSTRAINT uk_usr_login UNIQUE(login),
                       CONSTRAINT pk_user_id PRIMARY KEY(id)
);

CREATE TABLE lists
(
    id   int GENERATED ALWAYS AS IDENTITY NOT NULL,
    title text                             NOT NULL,
    CONSTRAINT pk_todo_id PRIMARY KEY (id)
);

CREATE TABLE userLists
(
    id int GENERATED ALWAYS AS IDENTITY NOT NULL,
    user_id int  NOT NULL,
    list_id int  NOT NULL,
    CONSTRAINT pk_ulists_id PRIMARY KEY (id),
    CONSTRAINT fk_ulists_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ,
    CONSTRAINT fk_ulists_list_id FOREIGN KEY (list_id) REFERENCES lists(id) ON DELETE CASCADE,
    CONSTRAINT uk_ulists_user_list UNIQUE (user_id, list_id)
);

CREATE TABLE cards
(
    id int GENERATED ALWAYS AS IDENTITY NOT NULL,
    front text                             NOT NULL,
    back text                               NOT NULL,
    imageLink text                              NOT NULL,
    voiceMessage text                       NOT NULL,
    CONSTRAINT pk_todo_items_id PRIMARY KEY (id)
);

CREATE TABLE listCards
(
    id int GENERATED ALWAYS AS IDENTITY NOT NULL,
    item_id int  NOT NULL,
    list_id int  NOT NULL,
    CONSTRAINT pk_lists_items_id PRIMARY KEY (id),
    CONSTRAINT fk_lists_items_item_id FOREIGN KEY (item_id) REFERENCES cards(id) ON DELETE CASCADE,
    CONSTRAINT fk_lists_items_list_id FOREIGN KEY (list_id) REFERENCES userLists(id) ON DELETE CASCADE
);