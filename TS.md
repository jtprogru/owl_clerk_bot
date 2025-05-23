# Technical specification

## Goals

Основная цель бота - принимать входящие от малознакомых людей и регистрировать предложения о работе от HR'ов.

## Functional requirements

Задачи бота:
- [ ] Принимать входящие сообщения;
- [ ] Регистрировать пользователя в БД;
- [ ] Собирать всю отправленную информацию от пользователя и передавать ее мне;

Дополнительные штуки:
- [ ] Веб-морда с различного рода информацией (кто/когда/сколько/откуда/etc);
- [ ] Возможность отправить конкретному пользователю сообщение прямо из веб-морды;
- [ ] Видеть контакты пользователя, которые он оставил;
- [ ] Отдельная регистрация HR'ов;


| func | input | output        |
|:---|:---|:--------------|
|SaveOrUpdateState| ctx context.Context, p IProfile, m IMessage| Answer, error |

| func | input                                     | output        |
|:---|:------------------------------------------|:--------------|
|SaveOrUpdateState| ctx context.Context, p Profile, m Message| answer, error |


IProfile
IMessage
Answer


| Id  | Next    | Answer  | Buttons | Handler |
|:----|:--------|:--------|:--------|:--------|
| 0   | 1       | Кто ты? |     |
| 1   | 1,2,3,4 | Ты HR?  | HR, NOT HR, CALC| ButtonSelect
| 2   | 2       | HR      | YES, NO | ButtonSelect
| 3   | 3       | NON HR  | YES, NO | ButtonSelect
| 4   | 5       | =+-*    | |Calculate
| 5   | 1       |         | |
