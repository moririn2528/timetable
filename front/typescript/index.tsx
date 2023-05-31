import React from "react";
import { createRoot } from 'react-dom/client';
import { ClassTimetableWithDate } from "./timetable_class";
import { TeacherTimetableWithDate } from "./timetable_teacher";
import { getClass, getTeachers, InputGraphType, TeacherType } from "./api";
import { GraphComponent } from "./class";

enum Menu {
    Classes = "クラス",
    ClassTimetable = "時間割",
    TeacherTimetable = "先生用時間割"
}

function ScreenWithMenu(props: { classes: InputGraphType, teachers: TeacherType[] }) {
    const [menu, setMenu] = React.useState(Menu.Classes);

    let componentToRender;
    switch (menu) {
        case Menu.Classes:
            componentToRender = <GraphComponent classes={props.classes} />;
            break;
        case Menu.ClassTimetable:
            componentToRender = <ClassTimetableWithDate classes={props.classes} date={new Date(2021, 3, 12)} />;
            break;
        case Menu.TeacherTimetable:
            componentToRender = <TeacherTimetableWithDate teachers={props.teachers} date={new Date(2021, 3, 13)} />;
            break;
        default:
            componentToRender = null;
            break;
    }
    return (
        <div style={{
            height: "100%",
        }}>
            <ul style={{
                display: "flex",
                width: "100%",
                height: "40px",
                backgroundColor: "dimgray",
                boxSizing: "border-box",
                zIndex: 1000,
            }}>
                {Object.values(Menu).map((m) => {
                    return <li style={{
                        listStyle: "none",
                        display: "block",
                        textDecoration: menu === m ? "underline" : "none",
                        color: "white",
                        margin: "8px 15px",
                    }} key={m} onClick={() => setMenu(m)}>{m}</li>
                })}
            </ul>
            {componentToRender}
        </div>
    )
}

window.addEventListener("load", async function (e) {
    const container = document.getElementById("root");
    if (!container) throw new Error("No root element")
    const classes = await getClass();
    const teachers = await getTeachers();
    const root = createRoot(container);
    root.render(<ScreenWithMenu classes={classes} teachers={teachers} />);
});