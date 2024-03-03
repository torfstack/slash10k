import discord
from discord.ext import commands
import os
import requests
import json

def base():
    return "https://true.torfstack.com/"

def getDebts():
    return base() + "api/debt"

def add10k(name:str):
    return base() + "api/debt/" + name + "/10000"

def sub10k(name:str):
    return base() + "api/debt/" + name + "/-10000"

intents = discord.Intents.default()
intents.message_content = True

client = commands.Bot(command_prefix='!', description="10k in die Gildenbank", intents=intents)
token = os.getenv('DISCORD_TOKEN')
if token is None:
    print("No token found")
    exit(1)

@client.event
async def on_ready():
    print(f'{client.user} has connected to Discord!')

@client.command(name='debts')
async def debts(ctx):
    debtsRes = requests.get(getDebts()).json()
    await ctx.send(''.join([debt["name"]+" "+debt["amount"]+"\n" for debt in debtsRes["debts"]]))
   
@client.command(name='+debt')
async def debts(ctx, name:str):
    addDebtRes = requests.post(add10k(name))
    debtsRes = requests.get(getDebts()).json()
    for debt in debtsRes["debts"]:
        if debt["name"] == name:
            await ctx.send(name+" "+debt["amount"])
            return
    
@client.command(name='-debt')
async def debts(ctx, name:str):
    addDebtRes = requests.post(sub10k(name))
    debtsRes = requests.get(getDebts()).json()
    for debt in debtsRes["debts"]:
        if debt["name"] == name:
            await ctx.send(name+" "+debt["amount"])
            return
        
client.run(token)
    