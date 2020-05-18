#!/usr/bin/python3

import argparse
import requests
import xml.etree.ElementTree as ET
import datetime
import json

class WeatherForecast:
	def __init__(self, latitude, longitude, start_datetime, end_datetime):
		self.latitude = latitude
		self.longitude = longitude
		self.start_datetime = start_datetime
		self.end_datetime = end_datetime
		self.product = "time-series"
		self.max_temperature = "maxt"
		self.min_temperature = "mint"
		self.probability_of_precip = "pop12"
		self.sky_cover = "sky"
		self.json_string_complete = ""

	def get_forecast(self):
		queryString = f"lat={self.latitude}&lon={self.longitude}&product={self.product}&begin={self.start_datetime}&end={self.end_datetime}&maxt={self.max_temperature}&mint={self.min_temperature}&pop12={self.probability_of_precip}&sky={self.sky_cover}"

		r = requests.get(f'https://graphical.weather.gov/xml/sample_products/browser_interface/ndfdXMLclient.php?{queryString}')
		self.forecast_info = r.text

		return self.forecast_info

	def get_total_cloud_coverage(self):
		jsonString = "\"cloudCover\": ["
		root = ET.fromstring(self.forecast_info)
		for item in root.findall('data/parameters/cloud-amount'):
			if item.attrib['type'] == 'total':
				time_layout = item.attrib['time-layout']
				(start_times,end_times) = self.get_time_layout(time_layout)
				current_index = 0
				for item2 in item.findall('value'):
					dictionaryData = {}
					dictionaryData['dateTime'] = f"{start_times[current_index]}"
					dictionaryData['percentage'] = item2.text
					if current_index > 0:
						jsonString = jsonString + ","
					jsonString = jsonString + json.dumps(dictionaryData, sort_keys=False, separators=(',', ': '))

					current_index = current_index + 1
		jsonString = jsonString + "]"
		if self.json_string_complete != "":
			self.json_string_complete = self.json_string_complete + ","

		self.json_string_complete = self.json_string_complete + jsonString


	def get_precip_probability(self):
		jsonString = "\"precipitationProbability\": ["
		root = ET.fromstring(self.forecast_info)
		for item in root.findall('data/parameters/probability-of-precipitation'):
			if item.attrib['type'] == '12 hour':
				time_layout = item.attrib['time-layout']
				(start_times,end_times) = self.get_time_layout(time_layout)
				current_index = 0
				for item2 in item.findall('value'):
					dictionaryData = {}
					dictionaryData['startDateTime'] = f"{start_times[current_index]}"
					dictionaryData['endDateTime'] = f"{end_times[current_index]}"
					dictionaryData['percentage'] = item2.text
					if current_index > 0:
						jsonString = jsonString + ","
					jsonString = jsonString + json.dumps(dictionaryData, sort_keys=False, separators=(',', ': '))

					current_index = current_index + 1
		jsonString = jsonString + "]"
		if self.json_string_complete != "":
			self.json_string_complete = self.json_string_complete + ","

		self.json_string_complete = self.json_string_complete + jsonString

	def get_max_temperature(self):
		jsonString = "\"maxTemperature\": ["
		root = ET.fromstring(self.forecast_info)
		for item in root.findall('data/parameters/temperature'):
			if item.attrib['type'] == 'maximum':
				time_layout = item.attrib['time-layout']
				(start_times,end_times) = self.get_time_layout(time_layout)
				current_index = 0
				for item2 in item.findall('value'):
					dictionaryData = {}
					dictionaryData['startDateTime'] = f"{start_times[current_index]}"
					dictionaryData['endDateTime'] = f"{end_times[current_index]}"
					dictionaryData['value'] = item2.text
					if current_index > 0:
						jsonString = jsonString + ","
					jsonString = jsonString + json.dumps(dictionaryData, sort_keys=False, separators=(',', ': '))

					current_index = current_index + 1
		jsonString = jsonString + "]"
		if self.json_string_complete != "":
			self.json_string_complete = self.json_string_complete + ","

		self.json_string_complete = self.json_string_complete + jsonString

	def get_min_temperature(self):
		jsonString = "\"minTemperature\": ["
		root = ET.fromstring(self.forecast_info)
		for item in root.findall('data/parameters/temperature'):
			if item.attrib['type'] == 'minimum':
				time_layout = item.attrib['time-layout']
				(start_times,end_times) = self.get_time_layout(time_layout)
				current_index = 0
				for item2 in item.findall('value'):
					dictionaryData = {}
					dictionaryData['startDateTime'] = f"{start_times[current_index]}"
					dictionaryData['endDateTime'] = f"{end_times[current_index]}"
					dictionaryData['value'] = item2.text
					if current_index > 0:
						jsonString = jsonString + ","
					jsonString = jsonString + json.dumps(dictionaryData, sort_keys=False, separators=(',', ': '))
					current_index = current_index + 1
		jsonString = jsonString + "]"
		if self.json_string_complete != "":
			self.json_string_complete = self.json_string_complete + ","

		self.json_string_complete = self.json_string_complete + jsonString

	def get_time_layout(self,layout_key):
		root = ET.fromstring(self.forecast_info)
		for item in root.findall('data/time-layout'):
			for item2 in item.findall('layout-key'):
				if item2.text == layout_key:
					start_times = []
					end_times = []
					for item3 in item.findall('start-valid-time'):
						start_times.append(item3.text)
					for item3 in item.findall('end-valid-time'):
						end_times.append(item3.text)
					return (start_times,end_times)
	
	def print_json(self):
		self.json_string_complete = "{" + self.json_string_complete + "}"
		print(json.dumps(json.loads(str(self.json_string_complete)), sort_keys=False, indent=4, separators=(',', ': ')))



parser = argparse.ArgumentParser()
parser.add_argument("--latitude", help="Latitude")
parser.add_argument("--longitude", help="Longitude")
parser.add_argument("--startDateTime", help="Starting date and time of forecast, e.g., 2017-07-29T00:00:00 (default to now)")
parser.add_argument("--endDateTime", help="Ending date and time of forecast, e.g., 2017-07-30T00:00:00 (default to max)")
args = parser.parse_args()

start_datetime = f"{datetime.date.today()}T00:00:00"
end_datetime = f"{datetime.date.today() + datetime.timedelta(days=4)}T23:59:59"

if args.startDateTime and args.endDateTime:
	cf = WeatherForecast(args.latitude, args.longitude, args.startDateTime, args.endDateTime)
else:
	# if no start and end date specified, then calculate a 5-day forecast
	cf = WeatherForecast(args.latitude, args.longitude, start_datetime, end_datetime)

forecast_data = cf.get_forecast()

cf.get_total_cloud_coverage()
cf.get_precip_probability()
cf.get_max_temperature()
cf.get_min_temperature()

cf.print_json()