# Chemical Mixing Data Synthesis
# Deep Learning  pH Prediction System

A smart web app for predicting pH changes during cream cheese fermentation using a hybrid **grey-box model**, combining **reaction kinetics** with **neural networks** (via TensorFlow). Inspired by **Guo et al. (2023)**.

---

##  Features

- Upload a CSV with fermentation time-series data
- Predict pH values using a kinetic + LSTM model
- Built with python backend and HTML/CSS frontend
- Integrated TensorFlow (Go/Python bindings)
- Frontend hosted with static HTML/CSS and Javascript
- Backend handles CSV parsing, kinetic simulation, and pH forecasting

---

## How to test product 
- Install Tensorflow, Flask and nake sure you have the latest version of python
- Clone github repository
- cd into GenerateData folder and run main.go folder which will generate a csv folder
    You can change various variables inthe data mixing process: initial starting biomass, the temprature, timesteps taken during cooking the variation periods .....
- cd into the lstme folder and train the lstm model with the csv models that you have genareted
- run the api
- Run the test development adress( the website is not upp and runing because of the incuring costs and low funding)
- Insert the One of the generated CSV files(lstm model  dosent use the ph column to make predictions)

## Probable use
Companies can integrate this method to save costs by knowing exact Ph value.
![UI Preview](chempic.png)
