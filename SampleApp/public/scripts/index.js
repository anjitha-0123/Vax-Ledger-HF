
// const addData = async (event) =>{
//     event.preventDefault();
//     const batchID = document.getElementById("batchID").value;
//     const manufacturer = document.getElementById("Manufacturer").value;
//     const vaccineType =document.getElementById("vaccineType").value;
//     const manufactureDate =document.getElementById("manufactureDate").value;
//     const expiryDate =document.getElementById("expiryDate").value;
//     // const minTemp =document.getElementById("minTemp").value;
//     // const maxTemp =document.getElementById("maxTemp").value;
//     const minTemp = parseInt(document.getElementById("minTemp").value);
//     const maxTemp = parseInt(document.getElementById("maxTemp").value);

    
//     const VaccineData = {
//         docType: "vaccineBatch",
//         batchID: batchID,
//         manufacturer:manufacturer,
//         vaccineType:vaccineType,
//         manufactureDate:manufactureDate,
//         expiryDate:expiryDate,
//         minTemp:minTemp,
//         maxTemp:maxTemp

//     };

//     console.log("Sending VaccineData:", VaccineData);
//     console.log("Type of minTemp:", typeof minTemp);
//     console.log("Type of maxTemp:", typeof maxTemp);

//     if (
//         batchID.length ==0 ||
//         manufacturer.length == 0 ||
//         vaccineType.length == 0 || 
//         manufactureDate.length == 0 || 
//         expiryDate.length == 0 || 
//         // minTemp.length == 0 || 
//         // maxTemp.length == 0 
//         isNaN(minTemp) ||
//         isNaN(maxTemp) 
//     ) {
//         alert("Please enter the data properly")
//     }
//    else {
//     try{
//         const response =await fetch("/api/batch", {
//             method: "POST",
//             headers: {
//                 'Content-Type':"application/json",
//             },
//             body:JSON.stringify(VaccineData)
//         });
//         console.log("RESPONSE :",response)
//         const data =await response.json()
//         console.log("DATA:",data);
//         return alert ("Vaccine Batch Created")

//     }catch(err){
//         alert ("Error");
//         console.log(err)
//     }
//    }
// }


// const readData =async (event)=>{
//     event.preventDefault();
//     const batchID = document.getElementById("batchidinput").value;
//     if (batchID.length==0){
//         alert ("Please Enter the valid ID")
//     }
//     else{
//         try {
//             const response =await fetch (`/api/batch/${batchID}`);
//             let responseData = await response.json();
//             console.log("response data:",responseData);
//             document.getElementById("resultBox").innerText = JSON.stringify(responseData, null, 2);
//             alert(JSON.stringify(responseData));

//         }
//         catch(err){
//             alert("Error");
//             console.log(err)
//         }
//     }
// };



