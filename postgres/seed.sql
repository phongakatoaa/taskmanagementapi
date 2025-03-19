-- password: password1
INSERT INTO auth.users (username, password, role)
VALUES ('Radahn', '$2a$10$C0GYbE0Kp2TESVvHW.v46utF.VybXCHm2OkGi35kwLTj1uurhKRae', 'EMPLOYEE');

-- password: password2
INSERT INTO auth.users (username, password, role)
VALUES ('Malenia', '$2a$10$pR84.oALjT/lDZWXJZUEvOWu7TPdPKw09jIjuTcPgxipNo43Ln08S', 'EMPLOYEE');

-- password: password3
INSERT INTO auth.users (username, password, role)
VALUES ('Tarnished', '$2a$10$ulG/fBM/kegfdh2g4n1uguHS3iSSDYX.GZy7bGJ8YiS13mhmUAHie', 'EMPLOYER');

INSERT INTO api.tasks (title, description, status, assigned_user_id, due_date)
VALUES ('Design database', 'design database', 'PENDING', 1, '2026-03-19T11:49:25+07:00');

INSERT INTO api.tasks (title, description, status, assigned_user_id, due_date)
VALUES ('Setup Docker', 'setup Docker', 'IN_PROGRESS', 2, '2026-03-20T11:49:25+07:00');

INSERT INTO api.tasks (title, description, status, assigned_user_id, due_date)
VALUES ('Build api', 'build api', 'COMPLETED', 1, '2026-03-21T11:49:25+07:00');