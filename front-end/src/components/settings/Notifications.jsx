import { useState } from "react";

function Notifications({emailNotification, email}) {
    const [selectedEmailNotification, setSelectedEmailNotification] = useState(emailNotification)

    return (
        <div className="p-4 w-full">
            <div className="mb-8">
                <p className="w-full text-3xl border-b py-2 mb-4">Email Notifications</p>
                <div>
                    <label>
                        <input type='checkbox' className="h-4 w-4 mr-2" checked={selectedEmailNotification?'checked':''} onChange={()=>setSelectedEmailNotification(!selectedEmailNotification)} defaultValue=''></input>
                        Enable Email Notifications to be sent to {email}
                    </label>
                </div>
            </div>
        </div>
    );
}

Notifications.defaultProps = {
    emailNotification: true,
    email: 'Username@mail.com'
}

export default Notifications;