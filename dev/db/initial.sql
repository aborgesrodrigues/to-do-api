    CREATE TABLE public."user" (
      id uuid NOT NULL,
      username varchar NOT NULL,
      "name" varchar NOT NULL,
      password varchar NOT NULL,
      CONSTRAINT user_pk PRIMARY KEY (id)
    );

    CREATE TABLE public.task (
      id uuid NOT NULL,
      user_id uuid NOT NULL,
      description varchar NOT NULL,
      state varchar NOT NULL,
      CONSTRAINT task_pk PRIMARY KEY (id)
    );


    -- public.task foreign keys

    ALTER TABLE public.task ADD CONSTRAINT task_fk FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE RESTRICT ON UPDATE RESTRICT;
    ALTER TABLE public."user" ADD CONSTRAINT user_unique UNIQUE (username);
