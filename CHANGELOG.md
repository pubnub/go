## [v5.0.0](https://github.com/pubnub/go/releases/tag/v5.0.0)
May-17-2021

#### Added
- Update Go SDK Metadata. 

#### Modified
- BREAKING CHANGE - IV used for encryption is now random by default. 
- BREAKING CHANGE - The SDK now suppoorts Go Modules. 

#### Fixed
- Presence event occupancy field parsed incorrectly. 

## [v4.10.0](https://github.com/pubnub/go/releases/tag/v4.10.0)
November-2-2020

#### Added
- Objects v2 PAM changes. 

#### Modified
- Fetch with 100 messages. 
- Include timetoken in send file. 
- Readme updates. 

#### Fixed
- Read Publish File Message retry count from config. 

## [v4.9.1](https://github.com/pubnub/go/releases/tag/v4.9.1)
October-1-2020

#### Fixed
- Fix for a deadlock on destroy. 
- Fetch response nil check. 

## [v4.9.0](https://github.com/pubnub/go/releases/tag/v4.9.0)
August-11-2020

#### Added
- History v3 with UUID and MessageType. 

#### Modified
- PNPublishMessage struct changes, PublishFile enhancements. 

## [v4.8.0](https://github.com/pubnub/go/releases/tag/v4.8.0)
July-24-2020

#### Added
- Files: Allows users to upload and share files.

#### Modified
- BREAKING CHANGE: EncryptString and DecryptString functions now accept a third param - bool, if true the IV is random and is sent along with the message. Default is false. 

#### Fixed
- BREAKING CHANGE: runes to string converion now returns string, this mostly affect the validation responses. This makes the SDK compatible with Go 1.15. 

## [v4.7.0](https://github.com/pubnub/go/releases/tag/v4.7.0)
June-10-2020

#### Added
- BREAKING CHANGE: This version does not support Objects v1 (beta). 

## [v4.6.6](https://github.com/pubnub/go/tree/v4.6.6)
  Apr-9-2020 

## [v4.6.5](https://github.com/pubnub/go/tree/v4.6.5)
  March-26-2020 

## [v4.6.4](https://github.com/pubnub/go/tree/v4.6.4)
  February-5-2020 

## [v4.6.3](https://github.com/pubnub/go/tree/v4.6.3)
  January-28-2020 

## [v4.6.2](https://github.com/pubnub/go/tree/v4.6.2)
  January-22-2020 

## [v4.6.1](https://github.com/pubnub/go/tree/v4.6.1)
  January-3-2020 

## [v4.6.0](https://github.com/pubnub/go/tree/v4.6.0)
  December-17-2019

## [v4.5.2](https://github.com/pubnub/go/tree/v4.5.2)
  November-27-2019

## [v4.5.1](https://github.com/pubnub/go/tree/v4.5.1)
  October-16-2019

## [v4.5.0](https://github.com/pubnub/go/tree/v4.5.0)
  October-8-2019

## [v4.4.0](https://github.com/pubnub/go/tree/v4.4.0)
  October-7-2019

## [v4.3.1](https://github.com/pubnub/go/tree/v4.3.1)
  October-2-2019

## [v4.3.0](https://github.com/pubnub/go/tree/v4.3.0)
  Septempber-23-2019

## [v4.2.7](https://github.com/pubnub/go/tree/v4.2.7)
  August-28-2019

- Add Objects method features

## [v4.0.0-beta.5](https://github.com/pubnub/go/tree/v4.0.0-beta.5)
  January-9-2018

- Add subscribe builder

## [v4.0.0-beta.4](https://github.com/pubnub/go/tree/v4.0.0-beta.4)
 December-20-2017

- Add Telemetry Manager

## [v4.0.0-beta.3](https://github.com/pubnub/go/tree/v4.0.0-beta.3)
 December-20-2017

- Add Destroy() method

## [v4.0.0-beta.2](https://github.com/pubnub/go/tree/v4.0.0-beta.2)
 November-7-2017

- Add reconnection manager
- Add HistoryDelete method
- Add demo app
- Rename channel group methods
- Fix signature generation

## [v4.0.0-beta](https://github.com/pubnub/go/tree/v4.0.0-beta)
 October-4-2017

- Beta release
