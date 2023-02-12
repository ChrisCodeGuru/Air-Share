import { useState } from "react";

function BillingPlans({planPurchased, planType, startDate, endDate, storageUsed, storageBought, autoRenewal}) {
    const [selectedAutoRenewal, setSelectedAutoRenewal] = useState(autoRenewal)

    return (
        <div className="p-4 w-full">
            <div className="mb-16">
                <p className="w-full text-3xl border-b py-2 mb-4  dark:border-white">Billings and Plans</p>
                {planPurchased ?
                    <div className="flex mb-8">
                        <div className="w-40 mr-8">
                            <p className="font-bold">Plan Purchased</p>
                            <div className="border border-black py-2 text-center h-16 dark:border-white">{planType}</div>
                        </div>
                        <div className="w-40 mr-8">
                            <p className="font-bold">Start Date</p>
                            <div className="border border-black py-2 text-center h-16 dark:border-white">{startDate}</div>
                        </div>
                        <div className="w-40 mr-8">
                            <p className="font-bold">End Date</p>
                            <div className="border border-black py-2 text-center h-16 dark:border-white">{endDate}</div>
                        </div>
                    </div>
                :
                <></>
                }
                <div>
                    <p className="font-bold">Date Usage</p>
                    <div className="flex w-[600px] max-w-[600px]">
                        <div className="w-[600px] h-2 rounded-md bg-gray-300"></div>
                        <div className="w-[300.33px] max-w-[600px] h-2 rounded-md bg-blue-500 absolute"></div>
                    </div>
                    <div className="flex w-[600px] justify-end">{storageUsed} GB/{storageBought} GB</div>
                </div>
                <div>
                    <label>
                        <input type='checkbox' className="h-4 w-4 mr-2" checked={selectedAutoRenewal?'checked':''} onChange={()=>setSelectedAutoRenewal(!selectedAutoRenewal)} defaultValue=''></input>
                        Enable Auto Renewal
                    </label>
                </div>
            </div>
        </div>
    );
}

BillingPlans.defaultProps = {
    planPurchased: true,
    planType: 'Annual 100 GB',
    startDate: '31 October 2022 24:22:30',
    endDate: '31 October 2023 24:22:30',
    storageUsed: 10.45,
    storageBought: 100,
    autoRenewal: false
}

export default BillingPlans;