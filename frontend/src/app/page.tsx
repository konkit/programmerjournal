'use client';

import TaskList from "@/components/taskList";
import Button from "@/components/button";
import {AddDay, DayOfWeek, Today, WallDate} from "@/lib/wall_date";
import {useState} from "react";

export default function Home() {

    const [todayDate, setTodayDate] = useState(Today())

    const changeDateForward = () => {
        setTodayDate((oldToday) => AddDay(oldToday, 1))
    }

    const changeDateBackward = () => {
        setTodayDate((oldToday) => AddDay(oldToday, -1));
    }

    return (
    <main className="flex min-h-screen flex-col items-center p-24">
        <div>
            <h1 className="text-3xl">
                <Button onClick={changeDateBackward}>&lt;&lt;</Button>
                {DayOfWeek(todayDate)}, {todayDate}
                <Button onClick={changeDateForward}>&gt;&gt;</Button>
            </h1>

            <TaskList tasks={[]} todayDate={todayDate}/>
        </div>
    </main>
    );
}
