import os
import re
import shutil
import sys
import uuid
import pyperclip
import win32con
from ctypes import *
import win32clipboard
from io import BytesIO
from PIL import ImageGrab
import win32clipboard as clip_board
import struct
from PIL import Image
from logHandler import logger

def set_clip_text(text=None, continueOnFailure=False, errMsg=''):
    """
    将文本添加到剪切板
    :param text:
    :param continueOnFailure:
    :param errMsg:
    :return:
    """
    logger.debug(f'将文本添加到剪切板--{text}')
    if isinstance(text, (tuple, list, dict, int, float, bool)):
        text = str(text)
    elif isinstance(text, str):
        pass
    else:
        logger.error("传入的文本格式错误")
        return False
    try:
        if sys.platform.lower() == "linux":
            pyperclip.copy(text)
            return True
        clip_board.OpenClipboard()
        clip_board.EmptyClipboard()
        clip_board.SetClipboardData(win32con.CF_UNICODETEXT, text)
        clip_board.CloseClipboard()
        return True
    except Exception as e:
        import pywintypes
        if type(e) == pywintypes.error:
            logger.error(f'文本框只支持填入字符串, {e}')
            return False
        logger.error('将文本添加到剪切板失败！')
        if continueOnFailure:
            return False
        else:
            errMsg = errMsg if errMsg else e
            raise Exception(errMsg)


def get_clip_text(continueOnFailure=False, errMsg=''):
    """
    获取剪切板的文本内容
    :param continueOnFailure:
    :param errMsg:
    :return:
    """

    try:
        if sys.platform.lower() == "linux":
            return pyperclip.paste()
        clip_board.OpenClipboard()
        text = clip_board.GetClipboardData(win32con.CF_UNICODETEXT)
        clip_board.CloseClipboard()
        return text
    except TypeError as e:
        logger.error('获取剪切板文本失败，原因是剪切板中不存在文本！')
        if continueOnFailure:
            if sys.platform.lower() == "win32":
                clip_board.CloseClipboard()
            return None
        else:
            errMsg = errMsg if errMsg else e
            raise Exception(errMsg)


def get_clip_image(mix_mode=False, continueOnFailure=False, errMsg=''):
    """
    获取剪切板图片路径
    :param mix_mode: 默认、混合模式
    :param continueOnFailure:
    :param errMsg:
    :return:
    """
    im = ImageGrab.grabclipboard()
    buffered = BytesIO()
    clip_image_arr = []
    try:
        if not mix_mode:
            if im is None:
                logger.error('剪切板可能同时存在图片和文字，请选择混合模式！')
                return []
            else:
                im.save(buffered, format='PNG')
                clip_img_path = extract_png_path()
                im.save(clip_img_path + "\\clip.jpg")
                clip_image_arr.append(clip_img_path + "\\clip.jpg")
                return clip_image_arr
        else:
            clip_board.OpenClipboard()
            clip_image_arr = []
            last_format = 0
            # 枚举format参数
            next_format = clip_board.EnumClipboardFormats(last_format)
            while next_format:
                # 枚举剪贴板上当前可用的数据格式-EnumClipboardFormats
                if next_format == 49341:
                    break
                clip_bytes = clip_board.GetClipboardData(next_format)
                # print(clip_board.GetClipboardData(next_format))
                if next_format == win32clipboard.CF_HDROP:
                    # 点击图标进行复制时，会是文件类型的
                    for line in clip_bytes:
                        if str(line).endswith(('.png', 'jpg')):
                            clip_image_arr.append(line)
                elif next_format in [win32clipboard.CF_DIB, win32clipboard.CF_DIBV5]:
                    # 此处可以优化长度和宽度
                    width = 588
                    height = 588
                    img = Image.frombytes('RGB', (width, height), clip_bytes)
                    clip_img_path = extract_png_path() + "\\clip" + str(uuid.uuid1()) + '.png'
                    img.save(clip_img_path)
                    logger.info('开始绘制')
                    clip_image_arr.append(clip_img_path)
                elif type(clip_bytes) == bytes:
                    # logger.info(f'{next_format}--')
                    if "QQRichEditFormat" in clip_bytes.decode("utf-8", "ignore") or \
                            "WeChatRichEditFormat" in clip_bytes.decode("utf-8", "ignore"):
                        qq_img_clip = clip_bytes.decode("utf-8", "ignore")
                        qq_img_clip = qq_img_clip.split("filepath=")[1]
                        qq_img_clip = qq_img_clip.split(" shortcut=")[0]
                        qq_img_clip = qq_img_clip.replace('"', '')
                        clip_image_arr.append(qq_img_clip)
                        # return qq_img_clip
                next_format = clip_board.EnumClipboardFormats(last_format)
                last_format = next_format
        clip_board.CloseClipboard()
        return clip_image_arr
    except Exception as e:
        logger.error('获取剪切板中的图片存放路径失败，原因是剪切板中没有获取到图片!')
        if continueOnFailure:
            return []
        else:
            errMsg = errMsg if errMsg else e
            raise Exception(errMsg)


def extract_png_path():
    img_dir = os.path.dirname(__file__)
    clip_img_path = img_dir + "\\clipImg"
    if not os.path.isdir(clip_img_path):
        os.makedirs(clip_img_path)
    return clip_img_path


class DropFiles(Structure):
    _fields_ = [
        ('pFiles', c_uint),
        ('x', c_long),
        ('y', c_long),
        ('fNC', c_int),
        ('fWide', c_bool)
    ]


def addFileToClipoard(fPath='', continueOnFailure=False, errMsg=''):
    """
    将文件添加到剪切板
    :param fPath: 文件路径，支持列表格式
    :param continueOnFailure: 失败后是否继续
    :param errMsg: 异常描述
    :return:
    """
    logger.debug(f'向剪切板写入文件--{fPath}')
    if isinstance(fPath, str):
        if not os.path.exists(fPath):
            logger.error(f'剪切板写入文件失败，指定的文件路径不存在：{fPath}')
            return False
        else:
            fPath = [fPath]
    elif isinstance(fPath, list):
        for f in fPath:
            if not os.path.exists(f):
                logger.error(f'剪切板写入文件失败，指定的文件路径不存在：{f}')
                return False
    else:
        logger.error('添加文件到剪切板失败，文件格式不正确！')

    
    try:
        if sys.platform.lower() == "linux":
            pyperclip.copy(fPath)
            return True
        win32clipboard.OpenClipboard()
        df = DropFiles()
        df.pFiles = sizeof(DropFiles)
        df.fWide = True

        fstr = '\0'.join(fPath) + '\0\0'
        data = fstr.encode('U16')[2:]

        win32clipboard.EmptyClipboard()
        win32clipboard.SetClipboardData(
            win32clipboard.CF_HDROP, bytes(df) + data
        )
        win32clipboard.CloseClipboard()
        return True
    except Exception as e:
        if sys.platform.lower() == "win32":
            win32clipboard.CloseClipboard()
        logger.error('添加文件到剪切板失败！')
        if continueOnFailure:
            return False
        else:
            errMsg = errMsg if errMsg else e
            raise Exception(errMsg)


def renameFile(fname: str = ''):
    if os.path.exists(fname):
        f = os.path.basename(fname)
        fTuple = os.path.splitext(f)
        ptn = '.*?-副本\((\d+)\)'
        res = re.findall(ptn, fTuple[0])
        if res:
            fixNum = 5 + len(res[0])
            fixStr = f'-副本({int(res[0]) + 1})'
            fnew = os.path.join(os.path.dirname(fname), fTuple[0][:-fixNum] + fixStr + fTuple[1])
        else:
            fnew = os.path.join(os.path.dirname(fname), fTuple[0] + '-副本(1)' + fTuple[1])
        return renameFile(fnew)
    else:
        return fname


def saveFromClipBoard(dirPath='', saveType="cover", continueOnFailure=False, errMsg=''):
    """
    保存剪切板文件
    :param dirPath: 文件夹路径
    :param saveType: 保存模式   cover-覆盖  cancel-放弃保存  rename-重命名
    :param continueOnFailure:
    :param errMsg:
    :return:
    """
    logger.debug('保存剪切板文件...')
   
    try:
        if sys.platform.lower() == "linux":
            files = pyperclip.paste()
        else:
            win32clipboard.OpenClipboard()
            files = win32clipboard.GetClipboardData(win32clipboard.CF_HDROP)   
            win32clipboard.CloseClipboard()         
        resFile = []
        for f in files:
            newfile = os.path.join(dirPath, os.path.basename(f))
            if saveType == "cover":
                # 覆盖模式
                shutil.copy(f, newfile)
                resFile.append(newfile)
            elif saveType == "cancel":
                # 放弃保存模式
                if os.path.exists(newfile):
                    resFile.append(None)
                else:
                    shutil.copy(f, newfile)
                    resFile.append(newfile)
            elif saveType == "rename":
                # 重命名模式
                newfile = renameFile(newfile)
                shutil.copy(f, newfile)
                resFile.append(newfile)
        
        return resFile

    except TypeError:
        logger.error('保存剪切板文件失败，原因是剪切板中没有文件！')
        if sys.platform.lower() == "win32":
            win32clipboard.CloseClipboard()
        return []
    except FileNotFoundError:
        logger.error('保存剪切板文件失败，目标文件夹不存在！')
        if sys.platform.lower() == "win32":        
            win32clipboard.CloseClipboard()
        return []
    except OSError:
        logger.error('保存剪切板文件失败，没有操作权限！')
        if sys.platform.lower() == "win32":        
            win32clipboard.CloseClipboard()
        return []
    except Exception as e:
        logger.error('保存剪切板文件失败，磁盘空间不足！')
        if sys.platform.lower() == "win32":        
            win32clipboard.CloseClipboard()
        if continueOnFailure:
            return []
        else:
            errMsg = errMsg if errMsg else e
            raise Exception(errMsg)


def emptyClipBoard(continueOnFailure=False, errMsg=''):
    """
    清空剪切板
    :param continueOnFailure:
    :param errMsg:
    :return:
    """
    try:
        if sys.platform.lower() == "linux":
            pyperclip.copy('')
            return True
        win32clipboard.OpenClipboard()
        win32clipboard.EmptyClipboard()
        win32clipboard.CloseClipboard()
        return True
    except Exception as e:
        logger.debug(f'清空当前剪切板内容失败，原因是：{e}')
        if continueOnFailure:
            return False
        else:
            errMsg = errMsg if errMsg else e
            raise Exception(errMsg)
